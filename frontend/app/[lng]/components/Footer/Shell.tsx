'use client'

import clsx from "clsx";
import { useState } from "react";
import { FaBars } from "react-icons/fa6";

export default function Shell({ lng, children }: { lng: string, children: React.ReactNode }) {
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  return (
    <div className='antialiased bg-gray-50 dark:bg-gray-900'>
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
              Hello This is Header
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
          className='overflow-y-auto py-5 px-3 h-full
        bg-white dark:bg-gray-800'>
          <div className='p-2'>
            Hello this is aside
          </div>
        </div>
      </aside>

      <main className="p-4 md:ml-64 h-auto pt-20">
        {children}
      </main>

    </div >
  )
}
