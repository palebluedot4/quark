import math
import sys
from collections.abc import Callable
from dataclasses import dataclass

import pytest

from mathx.abs import absolute, absolute_manual

type Absolute = Callable[[int], int]


@dataclass(frozen=True)
class Displacement:
    metres: float

    def __abs__(self) -> Displacement:
        return Displacement(abs(self.metres))


@pytest.mark.parametrize(
    "implementation",
    [
        pytest.param(absolute, id="stdlib"),
        pytest.param(absolute_manual, id="manual"),
    ],
)
@pytest.mark.parametrize(
    ("value", "expected"),
    [
        pytest.param(5, 5, id="positive"),
        pytest.param(-5, 5, id="negative"),
        pytest.param(0, 0, id="zero"),
        pytest.param(True, 1, id="true"),
        pytest.param(False, 0, id="false"),
        pytest.param(-(2**64), 2**64, id="beyond-64-bits"),
    ],
)
def test_absolute(
    implementation: Absolute,
    value: int,
    expected: int,
) -> None:
    actual = implementation(value)
    assert actual == expected
    assert type(actual) is int


@pytest.mark.parametrize(
    ("value", "expected"),
    [
        pytest.param(5.5, 5.5, id="positive"),
        pytest.param(-5.5, 5.5, id="negative"),
        pytest.param(math.inf, math.inf, id="positive-infinity"),
        pytest.param(-math.inf, math.inf, id="negative-infinity"),
        pytest.param(-sys.float_info.max, sys.float_info.max, id="largest"),
        pytest.param(-math.ulp(0.0), math.ulp(0.0), id="smallest-subnormal"),
    ],
)
def test_absolute_float(value: float, expected: float) -> None:
    actual = absolute(value)
    assert actual == expected
    assert type(actual) is float


def test_absolute_returns_positive_zero() -> None:
    actual = absolute(-0.0)
    assert math.copysign(1.0, actual) == 1.0


@pytest.mark.parametrize(
    "value",
    [
        pytest.param(math.nan, id="positive"),
        pytest.param(math.copysign(math.nan, -1), id="negative"),
    ],
)
def test_absolute_returns_positive_nan(value: float) -> None:
    actual = absolute(value)
    assert math.isnan(actual)
    assert math.copysign(1.0, actual) == 1.0


def test_absolute_complex_returns_magnitude() -> None:
    actual = absolute(3 + 4j)
    assert actual == 5.0
    assert type(actual) is float


def test_absolute_supports_dunder_abs() -> None:
    actual = absolute(Displacement(-3.5))
    assert actual == Displacement(3.5)
