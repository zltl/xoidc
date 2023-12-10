import Link from 'next/link'
import { useTranslation } from '@/app/i18n'



export default async function Page({ params: { lng } }: { params: { lng: string } }) {
  const { t } = await useTranslation(lng);

  return (
    <div> 
      this is main
      <div className='h-[3000px]'>
        abcdefg
      </div>
    </div>
  )
}
