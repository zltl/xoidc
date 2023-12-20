'use client'

import clsx from "clsx";
import Link from "next/link";
import { useState } from "react";
import { FaBars } from "react-icons/fa6";
import { useTranslation } from "@/app/i18n/client";
import { useCurrentPath } from "@/app/lib/client/useCurrentPath";


export default function Shell({ lng, children }: { lng: string, children: React.ReactNode }) {
  const { t } = useTranslation(lng, 'shell');
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const pathname = useCurrentPath();

  return (
    <div className='antialiased bg-gray-50 dark:bg-gray-900 h-full'>
      <nav className='bg-white border-b border-gray-200 px-4 py-2.5
       dark:bg-gray-800 dark:border-gray-700
       fixed left-0 right-0 top-0 z-50'>
        <div className='flex flex-wrap justify-between items-center'>
          <div className='flex justify-start items-center py-2'>
            <button className="transition-[width] w-6 md:w-0 overflow-x-hidden"
              onClick={() => setIsDrawerOpen(!isDrawerOpen)}>
              <FaBars />
            </button>
            <div>
              XOIDC
            </div>
          </div>
        </div>
      </nav>

      <aside
        className={clsx("fixed top-0 left-0 z-40 w-64 h-screen pt-14",
          "transition-transform  bg-white border-r",
          "border-gray-200 md:translate-x-0",
          "dark:bg-gray-800 dark:border-gray-700",
          {
            "-translate-x-full": !isDrawerOpen,
          })}
        aria-label="Sidenav"
        id="drawer-navigation"
      >
        <div
          className='overflow-y-auto py-5 h-full bg-white dark:bg-gray-800'>
          <Link
            href="/oidc-client"
            className={clsx('block py-2 px-3 text-base font-medium text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group',
              {
                'bg-gray-50 dark:bg-gray-700': pathname.startsWith(`/oidc-client`)
              })}>
            {t('OIDC Clients')}
          </Link>
          <Link
            href="/user-namespaces"
            className={clsx('block py-2 px-3 text-base font-medium text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group',
              {
                'bg-gray-50 dark:bg-gray-700': pathname.startsWith(`/user-namespaces`)
              })}>
            {t("User Namespaces")}
          </Link>
        </div>
      </aside>

      <main className="p-4 md:ml-64 pt-20">
        {children}
      </main>

    </div >
  )
}
