from PIL import Image
import io

def image_to_ascii(image_bytes, width=40):
    """Convert an image to ASCII art."""
    if not image_bytes:
        return ""
    
    try:
        img = Image.open(io.BytesIO(image_bytes))
        # Maintain aspect ratio
        aspect_ratio = img.height / img.width
        height = int(width * aspect_ratio * 0.5) # 0.5 because characters are taller than wide
        img = img.resize((width, height)).convert("L") # L mode is for grayscale
        
        chars = "@%#*+=-:. "
        pixels = img.getdata()
        ascii_str = ""
        for i, pixel in enumerate(pixels):
            ascii_str += chars[pixel // 26] # 256 / 10 roughly 25.6
            if (i + 1) % width == 0:
                ascii_str += "\n"
        return ascii_str
    except Exception as e:
        return f"Error converting image: {e}"
