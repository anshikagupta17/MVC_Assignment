"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";

export default function VillagePage() {
    console.log("VillagePage rendered");
    const [village, setVillage] = useState(null);
    const [loading, setLoading] = useState(true);

    async function fetchVillage() {
        try {
        const res = await apiFetch("/village");
        const data = await res.json();
        console.log("DATA FROM API:", data);

        if (!res.ok) {
            throw new Error(data.message || "Failed to fetch village");
        }

        setVillage(data);
        } catch (err) {
        console.error(err);
        alert("Failed to load village");
        } finally {
        setLoading(false);
        }
    }
    useEffect(() => {
        console.log("Village state changed:", village);
    }, [village]);

    useEffect(() => {
        fetchVillage();
    }, []);

    if (loading) {
        return <div className="p-6">Loading village...</div>;
    }

    if (!village) {
        console.log(village);
        return <div className="p-6">No village found</div>;
    }

    return (
        <div className="p-6 space-y-6">

            <h1 className="text-3xl font-bold">
            </h1>

            <div className="border p-4 space-y-2">
                <h2 className="text-xl font-semibold">Resources</h2>
                <p>Gold: {village.gold}</p>
                <p>Elixir: {village.elixir}</p>
                <p>Trophies: {village.trophies}</p>
            </div>

            <div className="border p-4">
                <h2 className="text-xl font-semibold">Townhall</h2>
                <p>Level: {village.townhall_level}</p>
            </div>
           
            <div className="border p-4">
            <h2 className="text-xl font-semibold mb-3">
                Buildings
            </h2>

            {village.buildings && village.buildings.length > 0 ? (
                <div className="space-y-2">
                {village.buildings.map((b) => (
                    <div key={b.id} className="border p-2">
                    <p>
                        Building ID: {b.building_id}
                    </p>
                    <p>
                        Level: {b.level}
                    </p>
                    <p>
                        Position: ({b.x}, {b.y})
                    </p>
                    <p>
                        Upgrading:{" "}
                        {b.upgrade_ends_at
                        ? "Yes"
                        : "No"}
                    </p>
                    </div>
                ))}
                </div>
            ) : (
                <p>No buildings found</p>
            )}
            </div>

                    </div>
    );
}