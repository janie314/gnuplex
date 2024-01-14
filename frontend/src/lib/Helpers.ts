class Helpers {
  public static fmtTime(rawtime: number): string {
    const secs = Math.floor(rawtime % 60);
    const mins = Math.floor(((rawtime - secs) % 3600) / 60);
    const hrs = Math.floor((rawtime - 60 * mins - secs) / 3600);
    const secs_str = secs.toString().padStart(2, "0");
    const mins_str = mins.toString().padStart(2, "0");
    const hrs_str = hrs === 0 ? "" : (hrs.toString() + ":");
    return `${hrs_str}${mins_str}:${secs_str}`;
  }
}

export { Helpers };
