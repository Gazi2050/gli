from PIL import Image

class ASCIIConverter:
    """
    Utility for converting image data into terminal-ready ASCII art.
    """
    
    ASCII_CHARS = ["@", "#", "S", "%", "?", "*", "+", ";", ":", ",", "."]

    def _resize_image(self, image, new_width=100):
        """Internal helper to scale images while maintaining aspect ratio."""
        width, height = image.size
        ratio = height / width / 1.65
        new_height = int(new_width * ratio)
        return image.resize((new_width, new_height))

    def _grayify(self, image):
        """Convert an RGB image to grayscale."""
        return image.convert("L")

    def _pixels_to_ascii(self, image):
        """Map grayscale pixel values to ASCII characters."""
        pixels = list(image.getdata())
        characters = "".join([self.ASCII_CHARS[pixel // 25] for pixel in pixels])
        return characters

    def convert(self, image_url, width=100):
        """
        Download and convert an image URL into a multi-line ASCII string.

        Args:
            image_url (str): The URL of the image to process.
            width (int, optional): The target width of the ASCII output. Defaults to 100.

        Returns:
            str: The formatted ASCII art string.
        """
        import requests
        from io import BytesIO

        try:
            response = requests.get(image_url)
            img = Image.open(BytesIO(response.content))
            
            img = self._resize_image(img, width)
            img = self._grayify(img)
            
            ascii_str = self._pixels_to_ascii(img)
            pixel_count = len(ascii_str)
            return "\n".join([ascii_str[index : (index + width)] for index in range(0, pixel_count, width)])
        except Exception as e:
            return f"Error converting image: {e}"
