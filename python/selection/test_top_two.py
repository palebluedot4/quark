import re
from collections.abc import Callable, Iterable

import pytest

from selection.top_two import top_two, top_two_manual

type TopTwo = Callable[[Iterable[int]], tuple[int, int]]


@pytest.mark.parametrize(
    "implementation",
    [
        pytest.param(top_two, id="stdlib"),
        pytest.param(top_two_manual, id="manual"),
    ],
)
class TestTopTwo:
    @pytest.mark.parametrize(
        ("values", "expected"),
        [
            pytest.param([1, 2], (2, 1), id="ascending_pair"),
            pytest.param([2, 1], (2, 1), id="descending_pair"),
            pytest.param([1, 2, 3, 4, 5], (5, 4), id="maximum_last"),
            pytest.param([5, 4, 3, 2, 1], (5, 4), id="maximum_first"),
            pytest.param([2, 3, 1], (3, 2), id="runner_up_before_maximum"),
            pytest.param([5, 3, 5], (5, 5), id="duplicate_maximum"),
            pytest.param([-3, -1, -2], (-1, -2), id="negative"),
        ],
    )
    def test_returns_the_top_two(
        self,
        implementation: TopTwo,
        values: list[int],
        expected: tuple[int, int],
    ) -> None:
        actual = implementation(values)
        assert actual == expected

    @pytest.mark.parametrize(
        ("values", "expected"),
        [
            pytest.param((3, 1, 2), (3, 2), id="tuple"),
            pytest.param(range(1, 4), (3, 2), id="range"),
            pytest.param({1, 3, 2}, (3, 2), id="set"),
        ],
    )
    def test_accepts_any_iterable(
        self,
        implementation: TopTwo,
        values: Iterable[int],
        expected: tuple[int, int],
    ) -> None:
        actual = implementation(values)
        assert actual == expected

    def test_consumes_an_iterator_once(self, implementation: TopTwo) -> None:
        actual = implementation(iter([1, 3, 2]))
        assert actual == (3, 2)

    def test_does_not_modify_the_input(self, implementation: TopTwo) -> None:
        values = [3, 1, 4, 1, 5]
        implementation(values)
        assert values == [3, 1, 4, 1, 5]

    @pytest.mark.parametrize(
        "values",
        [
            pytest.param([], id="empty"),
            pytest.param([1], id="single"),
        ],
    )
    def test_rejects_fewer_than_two_values(
        self,
        implementation: TopTwo,
        values: list[int],
    ) -> None:
        with pytest.raises(ValueError, match="not enough values"):
            implementation(values)


def test_top_two_manual_names_itself_in_the_error() -> None:
    message = "top_two_manual() not enough values (expected at least 2)"
    with pytest.raises(ValueError, match=rf"^{re.escape(message)}$"):
        top_two_manual([1])
