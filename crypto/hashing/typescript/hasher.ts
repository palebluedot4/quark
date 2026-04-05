export interface Hasher {
  hash(data: Uint8Array): Promise<Uint8Array>;
  verify(hash: Uint8Array, data: Uint8Array): Promise<boolean>;
}

export class SHA256Hasher implements Hasher {
  async hash(data: Uint8Array): Promise<Uint8Array> {
    const buffer = await crypto.subtle.digest("SHA-256", data as BufferSource);
    return new Uint8Array(buffer);
  }

  async verify(hash: Uint8Array, data: Uint8Array): Promise<boolean> {
    const computed = await this.hash(data);
    return timingSafeEqual(hash, computed);
  }
}

function timingSafeEqual(a: Uint8Array, b: Uint8Array): boolean {
  if (a.length !== b.length) {
    return false;
  }
  let mismatch = 0;
  for (let i = 0; i < a.length; i++) {
    mismatch |= a[i] ^ b[i];
  }
  return mismatch === 0;
}
