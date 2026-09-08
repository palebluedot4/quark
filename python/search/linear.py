from collections.abc import Sequence


def linear_search[T](values: Sequence[T], target: T) -> int:
    return values.index(target)


def linear_search_manual[T](values: Sequence[T], target: T) -> int:
    for i, value in enumerate(values):
        if value is target or value == target:
            return i
    raise ValueError("linear_search_manual() target not in sequence")
