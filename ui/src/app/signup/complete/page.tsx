
"use client";

import Image from "next/image";
import { Header } from '@/components/Header'

const gofundmeUrl = "https://www.gofundme.com/f/support-local-businesses-communities-with-mylocal";

export default function Complete() {
    return (
        <>
            <Header />
            <div className="flex flex-col px-6 py-12 lg:px-8">
                <div className="sm:mx-auto sm:w-full sm:max-w-sm">
                    <Image
                        width={110}
                        height={51}
                        className="mx-auto h-30 w-auto"
                        src="/myLocal.svg"
                        alt="myLocal Logo"
                        priority
                    />
                    <h2 className="mt-10 text-center text-xl/9 tracking-tight text-gray-900">
                        Thanks for subscribing!
                    </h2>
                    <div className="my-4 flex justify-center">

                        <button
                            onClick={() => window.open(gofundmeUrl, "_blank", "noopener,noreferrer")}
                            type="button"
                            className="rounded-md bg-yellow-600 px-4 py-2 text-sm font-semibold text-white shadow-xs hover:bg-yellow-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-yellow-600"
                            >
                            Donate Now
                        </button>
                    </div>
                </div>
            </div>
        </>
    );
}
