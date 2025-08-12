package main

import (
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})

	app.Post("/img-upload", func(c *fiber.Ctx) error {
		vips.Startup(nil)
		defer vips.Shutdown()

		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Error when upload!"})
		}

		tempPath := "./temp_" + file.Filename
		if err := c.SaveFile(file, tempPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
		}

		image1, err := vips.NewImageFromFile(tempPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to load image"})
		}
		defer image1.Close()
		image1.Resize(0.25, vips.KernelLanczos3)

		ep := vips.NewDefaultJPEGExportParams()
		ep.Quality = 62
		ep.StripMetadata = true
		image1Bytes, _, err := image1.Export(ep)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to export image"})
		}

		if err := os.WriteFile("/uploads/output.jpg", image1Bytes, 0644); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save output"})
		}

		return c.Status(200).JSON("OK")
	})

	log.Fatal(app.Listen(":7000"))
}
