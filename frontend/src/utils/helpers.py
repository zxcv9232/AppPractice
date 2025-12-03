def format_price(price: float) -> str:
    if price >= 1000:
        return f"${price:,.2f}"
    elif price >= 1:
        return f"${price:.2f}"
    else:
        return f"${price:.4f}"

def format_change(change: float) -> str:
    sign = "+" if change >= 0 else ""
    return f"{sign}{change:.2f}%"

def get_change_color(change: float) -> str:
    return "green" if change >= 0 else "red"

def format_volume(volume: float) -> str:
    if volume >= 1_000_000_000:
        return f"${volume/1_000_000_000:.2f}B"
    elif volume >= 1_000_000:
        return f"${volume/1_000_000:.2f}M"
    elif volume >= 1_000:
        return f"${volume/1_000:.2f}K"
    else:
        return f"${volume:.2f}"

