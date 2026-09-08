import re
from collections.abc import Callable, Sequence

import pytest

from search.linear import linear_search, linear_search_manual

type LinearSearch[T] = Callable[[Sequence[T], T], int]


@pytest.mark.parametrize(
    "implementation",
    [
        pytest.param(linear_search, id="stdlib"),
        pytest.param(linear_search_manual, id="manual"),
    ],
)
class TestLinearSearch:
    @pytest.mark.parametrize(
        ("values", "target", "expected"),
        [
            pytest.param([3, 1, 4, 1, 5], 3, 0, id="first"),
            pytest.param([3, 1, 4, 1, 5], 4, 2, id="middle"),
            pytest.param([3, 1, 4, 1, 5], 5, 4, id="last"),
            pytest.param([3, 1, 4, 1, 5], 1, 1, id="duplicate_target"),
            pytest.param([7], 7, 0, id="single"),
        ],
    )
    def test_returns_the_index(
        self,
        implementation: LinearSearch[int],
        values: Sequence[int],
        target: int,
        expected: int,
    ) -> None:
        actual = implementation(values, target)
        assert actual == expected

    @pytest.mark.parametrize(
        ("values", "target", "expected"),
        [
            pytest.param((3, 1, 4), 4, 2, id="tuple"),
            pytest.param(range(3, 6), 4, 1, id="range"),
        ],
    )
    def test_accepts_any_sequence(
        self,
        implementation: LinearSearch[int],
        values: Sequence[int],
        target: int,
        expected: int,
    ) -> None:
        actual = implementation(values, target)
        assert actual == expected

    def test_finds_a_string(self, implementation: LinearSearch[str]) -> None:
        actual = implementation(["a", "c", "b"], "b")
        assert actual == 2

    def test_matches_by_identity(self, implementation: LinearSearch[float]) -> None:
        nan = float("nan")
        actual = implementation([1.0, nan, 2.0], nan)
        assert actual == 1

    def test_matches_negative_zero_with_zero(
        self,
        implementation: LinearSearch[float],
    ) -> None:
        actual = implementation([-0.0], 0.0)
        assert actual == 0

    @pytest.mark.parametrize(
        "values",
        [
            pytest.param([3, 1, 4, 1, 5], id="absent"),
            pytest.param([], id="empty"),
        ],
    )
    def test_rejects_missing_target(
        self,
        implementation: LinearSearch[int],
        values: Sequence[int],
    ) -> None:
        with pytest.raises(ValueError, match="not in"):
            implementation(values, 2)


def test_linear_search_manual_names_itself_in_the_error() -> None:
    message = "linear_search_manual() target not in sequence"
    with pytest.raises(ValueError, match=rf"^{re.escape(message)}$"):
        linear_search_manual([3, 1, 4], 2)
