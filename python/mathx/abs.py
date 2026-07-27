from typing import SupportsAbs


def absolute[T](x: SupportsAbs[T]) -> T:
    return abs(x)


def absolute_manual(x: int) -> int:
    return -x if x < 0 else +x
