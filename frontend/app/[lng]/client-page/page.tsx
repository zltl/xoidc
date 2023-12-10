
import Link from 'next/link'
import { useTranslation } from '@/app/i18n';
import { Footer } from '@/app/[lng]/components/Footer/client'
import { useState } from 'react'

export default  async function Page({ params: { lng } }: { params: { lng: string } }) {
  const { t } = await useTranslation(lng, 'client-page');
  let counter = 0;
  const setCounter = (i:number) => {
    counter = i;
  }

  return (
    <>
    </>
  )
}
