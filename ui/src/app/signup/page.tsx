"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { Header } from '@/components/Header'

export default function Home() {
    const [email, setEmail] = useState("");
    const [name, setName] = useState("");
    const [newsletter, setNewsletter] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const router = useRouter();
    
    async function handleSubmit(e: FormEvent) {
        e.preventDefault();
        setLoading(true);
        setError("");

        const payload = {
            email,
            name,
            newsletter,
        };

        try {
            const res = await fetch(
                `/api/signup`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                },
            );
            if (!res.ok) throw new Error("Failed to sign up. Please try again.");
            localStorage.setItem("signup_email", email);
            router.push("/signup/validate");
        } catch (err: unknown) {
            const message =
                err instanceof Error ? err.message : "An unexpected error occurred.";
            setError(message);
            setLoading(false);
        }
    }

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
                    <h2 className="mt-1 text-center text-xl/9 tracking-tight text-gray-900">
                        Coming soon!
                    </h2>
                </div>

                <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                    <form onSubmit={handleSubmit} className="space-y-6">

                        {/* name */}
                        <div>
                            <label htmlFor="name" className="block text-sm/6 font-medium text-gray-900">
                                Name
                            </label>
                            <div className="mt-2">
                                <input
                                    id="name"
                                    type="text"
                                    required
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-yellow-600 sm:text-sm/6"
                                />
                            </div>
                        </div>

                        {/* email */}
                        <div>
                            <label htmlFor="email" className="block text-sm/6 font-medium text-gray-900">
                                Email address
                            </label>
                            <div className="mt-2">
                                <input
                                    id="email"
                                    type="email"
                                    required
                                    autoComplete="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-yellow-600 sm:text-sm/6"
                                />
                            </div>
                        </div>

                        {/* newsletter toggle */}
                        <div className="flex gap-3">
                            <div className="flex h-6 shrink-0 items-center">
                                <div className="group grid size-4 grid-cols-1">
                                    <input
                                        id="newsletter"
                                        type="checkbox"
                                        checked={newsletter}
                                        onChange={() => setNewsletter(!newsletter)}
                                        className="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-yellow-600 checked:bg-yellow-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-yellow-600"
                                    />
                                    <svg viewBox="0 0 14 14" className="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white" fill="none">
                                        <path d="M3 8L6 11L11 3.5" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round"
                                            className={newsletter ? "opacity-100" : "opacity-0"} />
                                    </svg>
                                </div>
                            </div>
                            <label htmlFor="newsletter" className="block text-sm/6 text-gray-900">
                                Subscribe to our newsletter
                            </label>
                        </div>

                        {/* submit */}
                        <div>
                            <button
                                type="submit"
                                disabled={loading}
                                className="flex w-full justify-center rounded-md bg-yellow-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-yellow-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-yellow-600"
                            >
                                {loading ? "Submitting…" : "Sign up"}
                            </button>
                        </div>

                        {error && <p className="text-sm text-red-600">{error}</p>}
                    </form>

                    <p className="mt-10 text-center text-sm/6 text-gray-500">
                        By clicking the Sign up button you agree to our<br />
                        <a href="https://mylocal.ing/terms/" target="_blank" rel="noopener noreferrer"
                            className="font-semibold text-yellow-600 hover:text-yellow-400">
                            Terms of Service
                        </a>{" "}
                        and{" "}
                        <a href="https://mylocal.ing/privacy/" target="_blank" rel="noopener noreferrer"
                            className="font-semibold text-yellow-600 hover:text-yellow-400">
                            Privacy Policy
                        </a>
                    </p>
                </div>
            </div>
        </>
    );
}
