'use client'

import Link from 'next/link'
import Image from 'next/image'

import { Container } from '@/components/Container'
import { NavLink } from '@/components/NavLink'

export function Header() {
    return (
        <header className="py-10">
            <Container>
                <nav className="relative z-50 flex justify-center">
                    <div className="flex items-center md:gap-x-12">
                        <Link href="https://mylocal.ing/#" aria-label="Home">
                            <Image
                                width={40}
                                height={40}
                                className="w-10 h-auto"
                                src="/myLocalFlower.svg"
                                alt="myLocal Logo Flower"
                                priority
                            />
                        </Link>
                        <div className="hidden md:flex md:gap-x-6">
                            <NavLink href="https://mylocal.ing/#tools">Tools</NavLink>
                            <NavLink href="https://mylocal.ing/#contribute">Contribute</NavLink>
                            <NavLink href="https://mylocal.ing/#pricing">Pricing</NavLink>
                            <NavLink href="https://mylocal.ing/#faq">Questions</NavLink>
                        </div>
                    </div>
                </nav>
            </Container>
        </header>
    )
}
