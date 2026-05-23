from typing import SupportsAbs


def absolute[T](x: SupportsAbs[T]) -> T:
    return abs(x)


def absolute_manual[T: (int, float)](x: T) -> T:
    return x if x >= 0 else -x
