export const convertUnixToDate = (
  unixTimestamp: number | undefined
): string | null => {
  if (!unixTimestamp) return null;
  const date = new Date(unixTimestamp * 1000);
  return date.toLocaleString();
};

export const TABLE_PAGE_SIZE = 5;
