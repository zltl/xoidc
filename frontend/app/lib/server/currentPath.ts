import { headers } from 'next/headers';

export const kXRPath = 'x-r-path';

export function currentPath() {
  const headersList = headers();
  const curpath = headersList.get(kXRPath) || '/';
  return curpath;
}
