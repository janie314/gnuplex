import type React from "react";
import { useCallback, useRef } from "react";

interface UseLongPressOptions {
  onShortClick: () => void;
  onLongPress: () => void;
  duration?: number;
}

let lastTouchedElement: HTMLElement | null = null;
let touchActive = false;
let touchTriggered = false;
let touchEndedRecently = false;
let touchResetTimeout: ReturnType<typeof setTimeout> | null = null;

function useLongPress({
  onShortClick,
  onLongPress,
  duration = 400,
}: UseLongPressOptions) {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isLongPressRef = useRef(false);
  const startX = useRef<number | null>(null);
  const startY = useRef<number | null>(null);
  const touchId = useRef<number | null>(null);
  const movedRef = useRef(false);
  const MOVE_THRESHOLD = 10;

  const handleMouseDown = useCallback(
    (e?: React.MouseEvent) => {
      isLongPressRef.current = false;

      timeoutRef.current = setTimeout(() => {
        isLongPressRef.current = true;
        onLongPress();
      }, duration);
    },
    [duration, onLongPress],
  );

  const handleMouseUp = useCallback(
    (e?: React.MouseEvent) => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      if (!isLongPressRef.current && !touchTriggered && !touchEndedRecently) {
        onShortClick();
      }

      isLongPressRef.current = false;
      touchTriggered = false;
      touchEndedRecently = false;
    },
    [onShortClick],
  );

  const handleMouseLeave = useCallback((e?: React.MouseEvent) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
    touchTriggered = false;
    touchEndedRecently = false;
  }, []);

  const handleTouchStart = useCallback(
    (e: React.TouchEvent) => {
      if (!e.touches || e.touches.length === 0) return;

      const target = e.currentTarget as HTMLElement;
      const isNetNewTouch = lastTouchedElement !== target || !touchActive;

      if (!isNetNewTouch) return;

      lastTouchedElement = target;
      touchActive = true;
      touchTriggered = true;
      touchEndedRecently = false;

      if (touchResetTimeout) {
        clearTimeout(touchResetTimeout);
        touchResetTimeout = null;
      }

      const t = e.touches[0] as Touch;
      touchId.current = t.identifier;
      startX.current = t.clientX;
      startY.current = t.clientY;
      movedRef.current = false;
      isLongPressRef.current = false;

      const touchDuration = Math.max(duration, 700);

      timeoutRef.current = setTimeout(() => {
        isLongPressRef.current = true;
        onLongPress();
      }, touchDuration);
    },
    [duration, onLongPress],
  );

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (!touchId.current) return;
    for (let i = 0; i < e.touches.length; i++) {
      const t = e.touches[i] as Touch;
      if (t.identifier === touchId.current) {
        const dx = Math.abs((startX.current ?? 0) - t.clientX);
        const dy = Math.abs((startY.current ?? 0) - t.clientY);
        if (dx > MOVE_THRESHOLD || dy > MOVE_THRESHOLD) {
          movedRef.current = true;
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
          }
        }
        break;
      }
    }
  }, []);

  const handleTouchEnd = useCallback(
    (e?: React.TouchEvent) => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      if (!isLongPressRef.current && !movedRef.current) {
        onShortClick();
      }

      isLongPressRef.current = false;
      touchId.current = null;
      startX.current = null;
      startY.current = null;
      movedRef.current = false;
      touchActive = false;
      touchEndedRecently = true;

      if (touchResetTimeout) {
        clearTimeout(touchResetTimeout);
      }
      touchResetTimeout = setTimeout(() => {
        touchTriggered = false;
        touchEndedRecently = false;
        touchResetTimeout = null;
      }, 500);
    },
    [onShortClick],
  );

  const handleTouchCancel = useCallback((e?: React.TouchEvent) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    if (touchResetTimeout) {
      clearTimeout(touchResetTimeout);
      touchResetTimeout = null;
    }
    isLongPressRef.current = false;
    touchId.current = null;
    startX.current = null;
    startY.current = null;
    movedRef.current = false;
    touchActive = false;
    touchTriggered = false;
    touchEndedRecently = false;
  }, []);

  return {
    onMouseDown: handleMouseDown,
    onMouseUp: handleMouseUp,
    onMouseLeave: handleMouseLeave,
    onTouchStart: handleTouchStart,
    onTouchMove: handleTouchMove,
    onTouchEnd: handleTouchEnd,
    onTouchCancel: handleTouchCancel,
  } as React.HTMLAttributes<HTMLInputElement>;
}

export { useLongPress };
