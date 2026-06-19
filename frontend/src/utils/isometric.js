const TILE_WIDTH = 64;
const TILE_HEIGHT = 32;
const ORIGIN_X = (50 * TILE_WIDTH) / 2;

export function tileToScreen(tileX, tileY) {
    const screenX = (tileX - tileY) * (TILE_WIDTH / 2) + ORIGIN_X;
    const screenY = (tileX + tileY) * (TILE_HEIGHT / 2);
    return { screenX, screenY };
}

export function screenToTile(screenX, screenY) {
    const adjustedX = screenX - ORIGIN_X;
    const tileX = (adjustedX / (TILE_WIDTH / 2) + screenY / (TILE_HEIGHT / 2)) / 2;
    const tileY = (screenY / (TILE_HEIGHT / 2) - adjustedX / (TILE_WIDTH / 2)) / 2;
    return {
        tileX: Math.round(tileX),
        tileY: Math.round(tileY),
    };
}

export { TILE_WIDTH, TILE_HEIGHT };