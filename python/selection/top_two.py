import heapq
from collections.abc import Iterable


def top_two(values: Iterable[int]) -> tuple[int, int]:
    first, second = heapq.nlargest(2, values)
    return first, second


def top_two_manual(values: Iterable[int]) -> tuple[int, int]:
    iterator = iter(values)
    try:
        first = next(iterator)
        second = next(iterator)
    except StopIteration:
        raise ValueError(
            "top_two_manual() not enough values (expected at least 2)"
        ) from None
    if first < second:
        first, second = second, first
    for value in iterator:
        if value > first:
            first, second = value, first
        elif value > second:
            second = value
    return first, second
