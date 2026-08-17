import { describe, expect, expectTypeOf, test } from "vitest";

import { abs } from "./abs.ts";

describe("abs", () => {
  describe("number", () => {
    test.each([
      { name: "positive", value: 5, expected: 5 },
      { name: "negative", value: -5, expected: 5 },
      { name: "zero", value: 0, expected: 0 },
      {
        name: "positive infinity",
        value: Number.POSITIVE_INFINITY,
        expected: Number.POSITIVE_INFINITY,
      },
      {
        name: "negative infinity",
        value: Number.NEGATIVE_INFINITY,
        expected: Number.POSITIVE_INFINITY,
      },
      { name: "largest", value: -Number.MAX_VALUE, expected: Number.MAX_VALUE },
      {
        name: "smallest subnormal",
        value: -Number.MIN_VALUE,
        expected: Number.MIN_VALUE,
      },
    ])("returns the absolute value for $name", ({ value, expected }) => {
      const actual = abs(value);
      expect(actual).toBe(expected);
    });

    test("returns positive zero", () => {
      const actual = abs(-0);
      expect(actual).toBe(0);
    });

    test("propagates NaN", () => {
      const actual = abs(Number.NaN);
      expect(actual).toBeNaN();
    });

    test("returns number", () => {
      const actual = abs(-5);
      expectTypeOf(actual).toEqualTypeOf<number>();
    });
  });

  describe("bigint", () => {
    test.each([
      { name: "positive", value: 5n, expected: 5n },
      { name: "negative", value: -5n, expected: 5n },
      { name: "zero", value: 0n, expected: 0n },
      { name: "beyond 64 bits", value: -(2n ** 64n), expected: 2n ** 64n },
    ])("returns the absolute value for $name", ({ value, expected }) => {
      const actual = abs(value);
      expect(actual).toBe(expected);
    });

    test("returns bigint", () => {
      const actual = abs(-5n);
      expectTypeOf(actual).toEqualTypeOf<bigint>();
    });
  });

  test("rejects the union of its overloads", () => {
    const value = 5 as number | bigint;
    // @ts-expect-error
    abs(value);
  });
});
