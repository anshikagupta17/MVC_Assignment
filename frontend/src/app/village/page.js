"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import Building from "@/components/building"
import { TROOP_TYPES } from "@/constants/troops";
import Link from "next/link";
import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";
import ProtectedRoute from "@/components/ProtectedRoute";
import Image from "next/image";
import { BUILDING_NAMES } from "@/constants/buildings";


function VillagePage() {
    const [village, setVillage] = useState(null);
    const [loading, setLoading] = useState(true);
    const [selectedBuilding, setSelectedBuilding] = useState(null);
    const [showShop, setShowShop] = useState(false);
    const [shopBuildings, setShopBuildings] = useState([]);
    const [buildingToPlace, setBuildingToPlace] = useState(null);
    const [showTroops, setShowTroops] = useState(false);
    const [troops, setTroops] = useState([]);
    const [selectedTroop, setSelectedTroop] = useState(null);
    const [trainQuantity, setTrainQuantity] = useState(1);
    const [errorMsg, setErrorMsg] = useState("");
    const [successMsg, setSuccessMsg] = useState("");


    const { logout } = useAuth();

    const router = useRouter();

    function handleLogout() {
        logout();
        router.push("/login");
    }

    function showError(rawMessage, fallbackMessage) {
        const displayMessage = fallbackMessage || rawMessage;
        setErrorMsg(displayMessage);
        setTimeout(() => setErrorMsg(""), 5000);
    }

   
    function showSuccess(message) {
        setSuccessMsg(message);
        setTimeout(() => setSuccessMsg(""), 3000);
    }
    

    async function fetchVillage() {
        try {
        const res = await apiFetch("/village");
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || "Failed to fetch village");
        }

        setVillage(data);
        if (selectedBuilding) {
            const updated = data.buildings?.find(b => b.id === selectedBuilding.id);
            if (updated) setSelectedBuilding(updated);
        }
        } catch (err) {
            console.error(err);
            showError(err.message,"Failed to load village");
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchVillage();
    }, []);

    useEffect(() => {
        const interval = setInterval(() => {
            fetchVillage();
        }, 10000);

        return () => clearInterval(interval);
    }, []);

    function handleSelectBuilding(building) {
       if (selectedBuilding?.id===building.id) {
        setSelectedBuilding(null);
        return;
       }
        setSelectedBuilding(building);
        if (building.building_id===10){
            openTroops();
        }
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

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            await fetchVillage();

            showSuccess("Upgrade started!");
        } catch (err) {
            console.error(err);
            showError(err.message,"Couldn't start the upgrade. Try again.");
        }
    }
    function handleGridClick(e) {
        const rect = e.currentTarget.getBoundingClientRect();
        let x = Math.floor((e.clientX - rect.left) / 20);
        let y = Math.floor((e.clientY - rect.top) / 20);

        if (buildingToPlace) {
            const sizeX = buildingToPlace.size_x || 1;
            const sizeY = buildingToPlace.size_y || 1;
            x = Math.max(0, Math.min(x, 49 - (sizeX - 1)));
            y = Math.max(0, Math.min(y, 49 - (sizeY - 1)));
            placeBuilding(x, y);
            return;
        }

        if (!selectedBuilding) return;

        const sizeX = selectedBuilding.size_x || 1;
        const sizeY = selectedBuilding.size_y || 1;
        x = Math.max(0, Math.min(x, 49 - (sizeX - 1)));
        y = Math.max(0, Math.min(y, 49 - (sizeY - 1)));
        moveBuilding(x, y);
    }

    async function moveBuilding(x, y) {
        if (!selectedBuilding) return;

        try {
            const res = await apiFetch("/village/buildings/move", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    building_instance_id: selectedBuilding.id,
                    x,
                    y,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }
            setSelectedBuilding(null)

            await fetchVillage();
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't move the building there. Try a different spot.");
        }
    }

    async function openShop() {
        try {
            const res = await apiFetch("/village/shop");
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to load shop");
            }

            setShopBuildings(data);
            setShowShop(true);
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't open the shop right now. Try again.");
        }
    }

    function handlePickBuilding(building) {
        setBuildingToPlace(building);
        setShowShop(false);
    }

    async function placeBuilding(x, y) {
        try {
            const res = await apiFetch("/village/buildings/build", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    building_id: buildingToPlace.building_id,
                    x,
                    y,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            setBuildingToPlace(null);
            await fetchVillage();
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't place the building there. Try a different spot.");
        }
    }

    async function collectResources() {
        try {
            const res = await apiFetch("/village/collect", {
                method: "POST",
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to collect resources");
            }

            await fetchVillage();
            showSuccess(`Collected ${data.Gold} gold and ${data.Elixir} elixir`);
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't collect resources right now. Try again.");
        }
    }

    async function openTroops() {
        try {
            const res = await apiFetch("/village/troops");
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to load troops");
            }

            const owned = data || [];
            const merged = TROOP_TYPES.map((type) => {
                const ownedTroop = owned.find((own) => own.troop_id === type.troop_id);
                if (ownedTroop) {
                    return { ...ownedTroop, unlock_level: type.unlock_level}
                }
                return ownedTroop || {
                    troop_id: type.troop_id,
                    name: type.name,
                    level: 0,
                    quantity: 0,
                    damage: null,
                    max_health: null,
                    unlock_level: type.unlock_level,
                };
            });

            setTroops(merged);
            setShowTroops(true);
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't load your troops right now. Try again.");
        }
    }

    async function trainTroop() {
        if (!selectedTroop) return;

        try {
            const res = await apiFetch("/village/troops/train", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    troop_id: selectedTroop.troop_id,
                    quantity: trainQuantity,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            await openTroops();
            await fetchVillage();
            showSuccess("Troops trained!");
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't train troops right now. Try again.");
        }
    }

    async function upgradeTroop() {
        if (!selectedTroop) return;

        try {
            const res = await apiFetch("/village/troops/upgrade", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    troop_id: selectedTroop.troop_id,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            await openTroops();
            showSuccess("Troop upgrade started!");
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't upgrade troops right now. Try again.");
        }
    }
    if (loading) {
        return <div className="p-6">Loading village...</div>;
    }

    if (!village) {
        console.log(village);
        return <div className="p-6">No village found</div>;
    }

    return (
        <div className="relative w-screen h-screen overflow-hidden bg-black flex items-center justify-center"
        >
            {errorMsg && (
                <div className="fixed top-6 left-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded shadow-lg z-50">
                    {errorMsg}
                </div>
            )}

            {successMsg && (
                <div className="fixed top-6 left-1/2 -translate-x-1/2 bg-green-600 text-white px-4 py-2 rounded shadow-lg z-50">
                    {successMsg}
                </div>
            )}

            <div
                className="relative overflow-hidden flex items-center justify-center"
                onClick={handleGridClick}
                style={{
                    width: "1020px",
                    height: "1020px",
                    backgroundImage: `
                        linear-gradient(rgba(255,255,255,0.18) 1px, transparent 1px),
                        linear-gradient(90deg, rgba(255,255,255,0.18) 1px, transparent 1px),
                        url('/assets/grass.jpeg')
                    `,
                    backgroundSize: `20px 20px, 20px 20px`,
                    backgroundColor: "#3F704D",
                    backgroundRepeat: "repeat",
                }}
            >
                <div className="relative w-full h-full">
                    {village.buildings?.map((b) => (
                        <Building
                            key={b.id}
                            building={b}
                            isSelected={selectedBuilding?.id === b.id}
                            onSelect={handleSelectBuilding}
                        />
                    ))}
                </div>
            </div>

            <div className="fixed top-4 left-4 z-20 flex gap-6 items-center">
                <div className="bg-black/70 text-white rounded-full px-4 py-2 flex gap-6 items-center">
                    <span>Gold: {village.gold}</span>
                    <span>Elixir: {village.elixir}</span>
                    <span>Trophies: {village.trophies}</span>
                </div>

                <button
                    onClick={handleLogout}
                    className="bg-black/70 rounded-full w-10 h-10 flex items-center justify-center"
                >
                    <Image
                        src="/assets/logout3.png"
                        alt="Logout"
                        width={20}
                        height={20}
                    />
                </button>
            </div>

            <div className="fixed top-4 left-1/2 -translate-x-1/2 z-20 bg-black/70 text-white rounded-full px-4 py-2">
                Townhall Lv {village.townhall_level}
            </div>

            <Link href="/battle" className="fixed top-4 right-4 z-20">
                <button className="bg-red-500 text-white rounded-full px-4 py-2">
                    Battle
                </button>
            </Link>

            <button
                onClick={openShop}
                className="fixed bottom-4 right-4 z-20 bg-blue-500 text-white rounded-full px-4 py-2"
            >
                Shop
            </button>

            {buildingToPlace && (
                <div className="fixed bottom-4 left-1/2 -translate-x-1/2 z-20 bg-black/70 text-white px-4 py-2 rounded">
                    Click on the map to place: {buildingToPlace.name}
                </div>
            )}

            {selectedBuilding && (
                <div className="fixed bottom-4 left-4 z-20 bg-black/80 text-white p-4 rounded flex gap-4 items-start">
                    <Image
                        src={`/assets/buildings/${selectedBuilding.building_id}.png`}
                        alt={BUILDING_NAMES[selectedBuilding.building_id]}
                        width={60}
                        height={60}
                    ></Image>

                    <h2 className="font-semibold">{BUILDING_NAMES[selectedBuilding.building_id]}</h2>
                    <p>ID: {selectedBuilding.id}</p>
                    <p>Level: {selectedBuilding.level}</p>
                    {selectedBuilding.upgrade_ends_at && (
                        <p className="text-red-400">Upgrading...</p>
                    )}

                    {(selectedBuilding.building_id === 6 || selectedBuilding.building_id === 7) && (
                        <button
                            onClick={collectResources}
                            className="border px-4 py-2 mt-2 bg-yellow-500 text-white"
                        >
                            Collect
                        </button>
                    )}

                    {(selectedBuilding.building_id === 1 && selectedBuilding.level >= 4) ||
                    (selectedBuilding.building_id !== 1 && selectedBuilding.level >= village.townhall_level) ? (
                        <button className="border px-4 py-2 mt-2 bg-gray-300 text-gray-500 cursor-not-allowed" disabled>
                            Max Level Reached
                        </button>
                    ) : (
                        <button
                            onClick={upgradeBuilding}
                            className="border px-4 py-2 mt-2 bg-blue-500 text-white"
                        >
                            Upgrade
                        </button>
                    )}
                </div>
            )}

            {showShop && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className=" text-[yellow] p-6 rounded max-w-md w-full max-h-[80vh] overflow-y-auto">
                        <h2 className="text-xl font-bold mb-4">Shop</h2>

                        {shopBuildings.length === 0 ? (
                            <p>No buildings available</p>
                        ) : (
                            shopBuildings.map((b) => (
                                <div
                                    key={b.building_id}
                                    className="border p-3 mb-2 flex justify-between items-center"
                                >
                                    <div>
                                        <p className="font-semibold">{b.name}</p>
                                        <p className="text-sm">
                                            Cost: {b.cost} {b.cost_type}
                                        </p>
                                        <p className="text-sm">
                                            Owned: {b.current_count}/{b.max_quantity}
                                        </p>
                                    </div>

                                    <button
                                        onClick={() => handlePickBuilding(b)}
                                        className="border px-3 py-1 bg-pink-500 text-yellow"
                                    >
                                        Build
                                    </button>
                                </div>
                            ))
                        )}

                        <button
                            onClick={() => setShowShop(false)}
                            className="border px-4 py-2 mt-4"
                        >
                            Close
                        </button>
                    </div>
                </div>
            )}

            {showTroops && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-[#5ECBC8] text-white p-6 rounded max-w-md w-full max-h-[80vh] overflow-y-auto">
                        <div className="flex justify-between items-center mb-6">
                            <h2 className="text-xl font-bold">Troops</h2>
                        
                            <button onClick={() => { setShowTroops(false); setSelectedTroop(null); }} className="border px-4 py-2 mt-4">
                                Close
                            </button>
                        </div>

                        {troops.length === 0 ? (
                            <p>No troops trained yet</p>
                        ) : (
                            troops.map((t) => {
                                const isLocked = t.unlock_level > village.townhall_level;

                                return(
                                    <div
                                        key={t.troop_id}
                                        onClick={() => {
                                            if (!isLocked) setSelectedTroop(t);
                                        }}
                                        className={`border border-gray-600 p-3 mb-2 rounded ${
                                            isLocked
                                                ? "bg-gray-500 opacity-50 cursor-not-allowed"
                                                : "bg-[#A43B76] text-white cursor-pointer"
                                        } ${
                                            selectedTroop?.troop_id === t.troop_id ? "ring-2 ring-yellow-300" : ""
                                        }`}
                                    >
                                        <img
                                            src={`/assets/troops/${t.troop_id}.png`}
                                            alt={t.name}
                                            className="w-12 h-12"
                                        ></img>
                                        <p className="font-semibold">{t.name}</p>
                                        <p className="text-sm">Level: {t.level || "-"}</p>
                                        <p className="text-sm">Quantity: {t.quantity}</p>
                                        <p className="text-sm">Damage: {t.damage ?? "-"}</p>
                                        <p className="text-sm">Max Health: {t.max_health ?? "-"}</p>
                                        {isLocked && (
                                            <p className="text-xs mt-1">Unlocks at Townhall {t.unlock_level}</p>
                                        )}
                                    </div>
                                );
                            })
                        )}

                        {selectedTroop && (
                            <div className="border-t pt-4 mt-4">
                                <label className="block mb-2">
                                    Quantity to train:
                                    <input
                                        type="number"
                                        min="1"
                                        value={trainQuantity}
                                        onChange={(e) => setTrainQuantity(Number(e.target.value))}
                                        className="border ml-2 px-2 py-1 w-20 text-white"
                                    />
                                </label>

                                <div className="flex gap-2">
                                    <button onClick={trainTroop} className="border px-4 py-2 bg-green-500 text-white">
                                        Train
                                    </button>
                                    <button onClick={upgradeTroop} className="border px-4 py-2 bg-blue-500 text-white">
                                        Upgrade
                                    </button>
                                </div>
                            </div>
                        )}

                        
                    </div>
                </div>
            )}
        </div>
    );

}

export default function VillagePageWrapper() {
    return (
        <ProtectedRoute>
            <VillagePage />
        </ProtectedRoute>
    );
}