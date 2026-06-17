"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import { BUILDING_NAMES } from "@/constants/buildings";

export default function VillagePage() {
    const [village, setVillage] = useState(null);
    const [loading, setLoading] = useState(true);
    const [selectedBuilding, setSelectedBuilding] = useState(null);
    const [now, setNow] = useState(Date.now());
    

    async function fetchVillage() {
        try {
        const res = await apiFetch("/village");
        const data = await res.json();

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

    useEffect(() => {
        const interval = setInterval(() => {
            setNow(Date.now());
        }, 1000);

        return () => clearInterval(interval);
    }, []);

    if (loading) {
        return <div className="p-6">Loading village...</div>;
    }

    if (!village) {
        console.log(village);
        return <div className="p-6">No village found</div>;
    }

        async function upgradeBuilding() {
        if (!selectedBuilding) return;

        try {
            const res = await apiFetch("/village/buildings/upgrade", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    building_instance_id: selectedBuilding.id,
                }),
            });
            console.log(res.status);

            const text = await res.text();
            console.log(text);

            if (!res.ok) {
                throw new Error(text);
            }

            await fetchVillage();

            alert("Upgrade started");
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    }

    function getRemainingTime(endTime) {
        const diff = new Date(endTime) - new Date();

        if (diff <= 0) return null;

        const seconds = Math.floor(diff / 1000);
        const minutes = Math.floor(seconds / 60);
        const remSeconds = seconds % 60;

        return `${minutes}m ${remSeconds}s`;
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
                <h2 className="text-xl font-semibold mb-4">
                    Village Layout
                </h2>

                <div className="relative w-[600px] h-[600px] border bg-[#3F704D] overflow-hidden">

                    {village.buildings?.map((b) => (
                    <div
                        key={b.id}
                        onClick={() => setSelectedBuilding(b)}
                        className={`absolute w-16 h-16 border flex flex-col items-center justify-center text-xs text-center shadow cursor-pointer hover:scale-105 transition
                            ${b.upgrade_ends_at ? "opacity-60 border-red-400" : "bg-black"}
                        `}
                        style={{
                            left: `${b.x * 20}px`,
                            top: `${b.y * 20}px`,
                        }}
                    >
                        <div>
                        {BUILDING_NAMES[b.building_id]}
                        </div>
                        

                        <div>
                        <div>
                            Lv {b.level}
                        </div>

                        {b.upgrade_ends_at && (
                            <div className="text-red-500">
                                {getRemainingTime(b.upgrade_ends_at)}
                            </div>
                        )}
                        </div>
                    </div>
                    ))}

                </div>
                {selectedBuilding && (
                    <div className="border p-4 mt-4 space-y-2">
                        <h2 className="text-xl font-semibold">
                            Building Details
                        </h2>

                        <p>
                            Name: {BUILDING_NAMES[selectedBuilding.building_id]}
                        </p>

                        <p>
                            Level: {selectedBuilding.level}
                        </p>

                        <p>
                            Building ID: {selectedBuilding.id}
                        </p>

                        <p>
                            Position:
                            ({selectedBuilding.x}, {selectedBuilding.y})
                        </p>

                        <p>
                            Upgrade Status:
                            {selectedBuilding.upgrade_ends_at
                                ? " Upgrading"
                                : " Idle"}
                        </p>

                        <button
                            onClick={upgradeBuilding}
                            className="border px-4 py-2"
                        >
                            Upgrade
                        </button>
                    </div>
                )}
            </div>

        </div>
    );



}