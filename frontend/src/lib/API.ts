interface PosResponse {
  pos: number;
  max_pos: number;
}

const common_options = {
  headers: {
    "Content-Type": "application/json",
  },
};

class APICall {
  /**
   * @returns GNUPlex version string.
   */
  public static async version(): Promise<string> {
    return await fetch(
      "/api/version",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @returns Whether the video player is paused.
   */
  public static async paused(): Promise<boolean | null> {
    return await fetch(
      "/api/paused",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * Toggles the video's play/pause status.
   *
   * @returns Whether the video player is paused following the toggle operation.
   */
  public static async toggle(): Promise<boolean> {
    return await fetch(
      "/api/toggle",
      { method: "POST", ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @returns The video player's position (seconds).
   */
  public static async pos(): Promise<PosResponse | null> {
    return await fetch(
      "/api/pos",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * Sets the video's position, either with an absolute position or an increment.
   *
   * @param pos Position to set, either an increment or an absolute position.
   * @param inc Whether `pos` is an increment or absolute position.
   * @returns The video player's position (seconds) following the toggle inc operation.
   */
  public static async setPos(
    pos: number,
    inc: boolean,
  ): Promise<number | null> {
    return await fetch(
      "/api/pos",
      { method: "POST", ...common_options, body: JSON.stringify({ inc, pos }) },
    ).then((res) => res.json());
  }

  /**
   * @returns The video's volume (percentage, 0 - 100+).
   */
  public static async vol(): Promise<number | null> {
    return await fetch(
      "/api/vol",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @param vol The video's volume (percentage, 0 - 100+).
   */
  public static async setVol(vol: number) {
    return await fetch(
      "/api/vol",
      { method: "POST", ...common_options, body: JSON.stringify({ vol }) },
    ).then((res) => res.json());
  }

  public static async media(): Promise<string | null> {
    return await fetch(
      "/api/media",
      { ...common_options },
    ).then((res) => res.json());
  }

  public static async setMedia(mediafile: string) {
    return await fetch(
      "/api/media",
      {
        method: "POST",
        ...common_options,
        body: JSON.stringify({ media: mediafile }),
      },
    );
  }

  public static async getOriginMediadirs() {
    return await fetch(
      `/api/mediadirs`,
    ).then((res) => res.json()).then((data: string[]) => {
      return data.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
    });
  }

  public static async setOriginMediadirs(mediadirs: string[]) {
    return await fetch(
      `/api/mediadirs?mediadirs=${encodeURI(JSON.stringify(mediadirs))}`,
      { method: "POST" },
    );
  }

  public static async getOriginFileExts() {
    return await fetch(
      `/api/file_exts`,
    ).then((res) => res.json()).then((data: string[]) => {
      return data.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
    });
  }

  public static async setOriginFileExts(file_exts: string[]) {
    return await fetch(
      `/api/file_exts?file_exts=${encodeURI(JSON.stringify(file_exts))}`,
      { method: "POST" },
    );
  }

  public static async getOriginMediafiles() {
    return await fetch(`/api/medialist`).then((res) => res.json());
  }

  public static async refreshOriginMediafiles() {
    return await fetch(`/api/medialist`, { method: "POST" });
  }

  public static async getOriginLast25() {
    return await fetch(`/api/last25`).then((res) => res.json());
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export { type PosResponse };
export { APICall };
