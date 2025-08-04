import Pocketbase from 'pocketbase';

const pb = new Pocketbase('http://127.0.0.1:8090');

export default pb;

function formatUTC(date: Date): string {
    const iso = date.toISOString();
    return iso.replace('T', ' ').replace('Z', 'Z'); 
}

export function getUTCNowMinus(duration: { hours?: number, days?: number }): string {
    const now = new Date();
    const ms = now.getTime()
        - (duration.hours || 0) * 60 * 60 * 1000
        - (duration.days || 0) * 24 * 60 * 60 * 1000;

    const past = new Date(ms);
    return formatUTC(past);
}