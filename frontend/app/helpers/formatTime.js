export default function formatTime(ms) {
    let seconds = Math.floor((ms / 1000) % 60);
    let minutes = Math.floor((ms / (1000 * 60)) % 60);
    let hours = Math.floor((ms / (1000 * 60 * 60)) % 24);
    let days = Math.floor(ms / (1000 * 60 * 60 * 24));

    seconds = String(seconds).padStart(2, '0');
    minutes = String(minutes).padStart(2, '0');
    hours = String(hours).padStart(2, '0');
    days = String(days).padStart(2, '0');

    let timeString = `${minutes}:${seconds}`;
    if (hours !== "00") timeString = `${hours}:${timeString}`;
    if (days !== "00") timeString = `${days}:${timeString}`;

    return timeString;
}