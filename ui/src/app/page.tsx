"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { Header } from '@/components/Header'

export default function Home() {
    const [email, setEmail] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const router = useRouter();
    
    async function handleSubmit(e: FormEvent) {
        e.preventDefault();
        setLoading(true);
        setError("");

        const payload = {
            email,
        };

        try {
            const res = await fetch(
                `/api`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                },
            );
            if (!res.ok) throw new Error("Failed to sign in. Please try again.");
            localStorage.setItem("signin_email", email);
            router.push("/validate");
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
                </div>

                <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                    <form onSubmit={handleSubmit} className="space-y-6">

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

                        {/* submit */}
                        <div>
                            <button
                                type="submit"
                                disabled={loading}
                                className="flex w-full justify-center rounded-md bg-yellow-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-yellow-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-yellow-600"
                            >
                                {loading ? "Submitting…" : "Sign in"}
                            </button>
                        </div>

                        {error && <p className="text-sm text-red-600">{error}</p>}
                    </form>

                    <p className="mt-10 text-center text-sm/6 text-gray-500">
                        By clicking the Sign in button you agree to our<br />
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
