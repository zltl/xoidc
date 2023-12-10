import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';
import { kXRPath } from './app/lib/server/currentPath';

acceptLanguage.languages(languages);

export const config = {
  // matcher: '/:lng*'
  matcher: [
    '/((?!api|_next/static|_next/image|assets|favicon.ico|sw.js).*)', 
    '/',
  ],
}

export function middleware(req:NextRequest) {
  const basepath = process.env.NEXT_PUBLIC_BASEPATH || '';
  let pathname = req.nextUrl.pathname;
  if (req.nextUrl.pathname.startsWith(basepath)) {
    pathname = req.nextUrl.pathname.substring(basepath.length);
  }

  let lng
  if (req.cookies.has(cookieName)) lng = acceptLanguage.get(req.cookies.get(cookieName)?.value)
  if (!lng) lng = acceptLanguage.get(req.headers.get('Accept-Language'))
  if (!lng) lng = fallbackLng


  let xrpath = pathname;
  languages.forEach(loc => {
    xrpath = xrpath.replace(`/${loc}`, '');
  });
  const h = new Headers(req.headers);
  h.set(kXRPath, xrpath);

  // Redirect if lng in path is not supported
  if (
    !languages.some(loc => pathname.startsWith(`/${loc}`)) &&
    !pathname.startsWith('/_next')
  ) {
    return NextResponse.redirect(new URL(`${basepath}/${lng}/${pathname}`, req.url));
  }

  if (req.headers.has('referer')) {
    const refererUrl = new URL(req.headers.get('referer')!)
    const lngInReferer = languages.find((l) => refererUrl.pathname.startsWith(`${basepath}/${l}`))
    const response = NextResponse.next({
      request: {
        ...req,
        headers: h,
      }
    })
    if (lngInReferer) response.cookies.set(cookieName, lngInReferer)
    return response;
  }

  return NextResponse.next({
    request: {
      ...req,
      headers: h,
    }
  });
}
