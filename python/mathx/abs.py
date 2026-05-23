from typing import SupportsAbs


def absolute[T](x: SupportsAbs[T]) -> T:
    return abs(x)
