"use client";

import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function ProtectedRoute({ children }) {
    const { isAuthenticated, authLoading } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!authLoading && !isAuthenticated) {
            router.push("/login");
        }
    }, [authLoading, isAuthenticated]);

    if (authLoading) {
        return <div className="p-6">Loading...</div>;
    }

    if (!isAuthenticated) {
        return <div className="p-6">Redirecting to login...</div>;
    }

    return children;
}