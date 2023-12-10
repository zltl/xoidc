'use client';

import { languages } from "@/app/i18n/settings";
import { usePathname } from "next/navigation";

export function useCurrentPath() {
  let curpath = usePathname();
  // remove basepath
  const basepath = process.env.NEXT_PUBLIC_BASEPATH || '';
  if (curpath.startsWith(basepath)) {
    curpath = curpath.substring(basepath.length);
  }
  // remove all languages from path
  languages.forEach(loc => {
    curpath = curpath.replace(`/${loc}`, '');
  });
  return curpath;
}
