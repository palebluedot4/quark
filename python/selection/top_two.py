import heapq
from collections.abc import Iterable


def top_two(values: Iterable[int]) -> tuple[int, int]:
    first, second = heapq.nlargest(2, values)
    return first, second
