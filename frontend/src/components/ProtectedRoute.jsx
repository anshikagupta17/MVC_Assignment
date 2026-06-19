"use client";

import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function ProtectedRoute({ children }) {
    const { isAuthenticated } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!isAuthenticated) {
            router.push("/login");
        }
    }, [isAuthenticated]);

    if (!isAuthenticated) {
        return <div className="p-6">Redirecting to login...</div>;
    }

    return children;
}