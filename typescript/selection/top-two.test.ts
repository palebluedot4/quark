import { describe, expect, expectTypeOf, test } from "vitest";

import { type TopTwo, topTwo } from "./top-two.ts";

describe("topTwo", () => {
  test.each([
    { name: "ascending pair", values: [1, 2], expected: [2, 1] },
    { name: "descending pair", values: [2, 1], expected: [2, 1] },
    { name: "maximum last", values: [1, 2, 3, 4, 5], expected: [5, 4] },
    { name: "maximum first", values: [5, 4, 3, 2, 1], expected: [5, 4] },
    { name: "runner-up before maximum", values: [2, 3, 1], expected: [3, 2] },
    { name: "duplicate maximum", values: [5, 3, 5], expected: [5, 5] },
    { name: "negative", values: [-3, -1, -2], expected: [-1, -2] },
    {
      name: "negative infinity",
      values: [Number.NEGATIVE_INFINITY, Number.NEGATIVE_INFINITY],
      expected: [Number.NEGATIVE_INFINITY, Number.NEGATIVE_INFINITY],
    },
  ])("returns the top two for $name values", ({ values, expected }) => {
    const actual = topTwo(values);
    expect(actual).toEqual(expected);
  });

  test.each([
    { name: "leading NaN", values: [Number.NaN, 1], expected: [1, Number.NaN] },
    {
      name: "trailing NaN",
      values: [1, Number.NaN],
      expected: [1, Number.NaN],
    },
    {
      name: "NaN between numbers",
      values: [1, Number.NaN, 2],
      expected: [2, 1],
    },
    {
      name: "NaN after numbers",
      values: [1, 2, Number.NaN],
      expected: [2, 1],
    },
    {
      name: "single non-NaN",
      values: [Number.NaN, Number.NaN, 1],
      expected: [1, Number.NaN],
    },
  ])("orders NaN last for $name", ({ values, expected }) => {
    const actual = topTwo(values);
    expect(actual).toEqual(expected);
  });

  test.each([
    { name: "empty", values: [] },
    { name: "single", values: [1] },
  ])("returns undefined for $name values", ({ values }) => {
    const actual = topTwo(values);
    expect(actual).toBeUndefined();
  });

  test("does not modify the input", () => {
    const values = [3, 1, 4, 1, 5];
    topTwo(values);
    expect(values).toEqual([3, 1, 4, 1, 5]);
  });

  test("narrows the return type for a tuple of at least two", () => {
    const actual = topTwo([1, 3, 2] as const);
    expectTypeOf(actual).toEqualTypeOf<TopTwo>();
  });

  test("does not narrow the return type for a single-element tuple", () => {
    const actual = topTwo([1] as const);
    expectTypeOf(actual).toEqualTypeOf<TopTwo | undefined>();
  });

  test("includes undefined for an array that may be short", () => {
    const values: readonly number[] = [1, 3, 2];
    const actual = topTwo(values);
    expectTypeOf(actual).toEqualTypeOf<TopTwo | undefined>();
  });
});
