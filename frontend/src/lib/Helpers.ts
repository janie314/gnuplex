export interface TimeComponents {
  hours: number;
  minutes: number;
  seconds: number;
}

export function secondsToTimeComponents(totalSeconds: number): TimeComponents {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = Math.floor(totalSeconds % 60);
  return { hours, minutes, seconds };
}

export function timeComponentsToSeconds(time: TimeComponents): number {
  return time.hours * 3600 + time.minutes * 60 + time.seconds;
}
