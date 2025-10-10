function uuidToIntArray(uuidStr) {
  // 하이픈 제거
  const hex = uuidStr.replace(/-/g, "");
  if (hex.length !== 32) throw new Error("Invalid UUID format");

  // 상위/하위 64비트
  const msb = BigInt("0x" + hex.slice(0, 16));
  const lsb = BigInt("0x" + hex.slice(16));

  // int 배열로 분리
  const f = Number((msb >> 32n) & 0xffffffffn) | 0;
  const g = Number(msb & 0xffffffffn) | 0;
  const h = Number((lsb >> 32n) & 0xffffffffn) | 0;
  const i = Number(lsb & 0xffffffffn) | 0;

  return [f, g, h, i];
}
function newUUID() {
  let result = "";
  let myuuid = crypto.randomUUID();
  result += myuuid.toString() + "\n";
  result += uuidToIntArray(myuuid.toString());
  return result;
}
newUUID();
