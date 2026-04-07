// @ts-check

/**
 * @typedef {Object} Hasher
 * @property {(data: Uint8Array) => Promise<Uint8Array>} hash
 * @property {(hash: Uint8Array, data: Uint8Array) => Promise<boolean>} verify
 */

/**
 * @implements {Hasher}
 */
export class SHA256Hasher {
  /**
   * @param {Uint8Array} data
   * @returns {Promise<Uint8Array>}
   * @throws {TypeError}
   */
  async hash(data) {
    if (!(data instanceof Uint8Array)) {
      throw new TypeError(`data must be Uint8Array, got: ${typeof data}`);
    }
    const buffer = await crypto.subtle.digest(
      "SHA-256",
      /** @type {BufferSource} */ (data),
    );
    return new Uint8Array(buffer);
  }

  /**
   * @param {Uint8Array} hash
   * @param {Uint8Array} data
   * @returns {Promise<boolean>}
   * @throws {TypeError}
   */
  async verify(hash, data) {
    if (!(hash instanceof Uint8Array) || !(data instanceof Uint8Array)) {
      throw new TypeError("arguments must be Uint8Array");
    }
    const computed = await this.hash(data);
    return timingSafeEqual(hash, computed);
  }
}

/**
 * @param {Uint8Array} a
 * @param {Uint8Array} b
 * @returns {boolean}
 */
export function timingSafeEqual(a, b) {
  if (!(a instanceof Uint8Array) || !(b instanceof Uint8Array)) {
    throw new TypeError("arguments must be Uint8Array");
  }
  if (a.length !== b.length) {
    return false;
  }
  let mismatch = 0;
  for (let i = 0; i < a.length; i++) {
    mismatch |= a[i] ^ b[i];
  }
  return mismatch === 0;
}
