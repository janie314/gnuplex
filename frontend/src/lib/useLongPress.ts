import { useCallback, useRef } from "react";

interface UseLongPressOptions {
  onShortClick: () => void;
  onLongPress: () => void;
  duration?: number;
}

function useLongPress({
  onShortClick,
  onLongPress,
  duration = 400,
}: UseLongPressOptions) {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isLongPressRef = useRef(false);

  const handleMouseDown = useCallback(() => {
    isLongPressRef.current = false;

    timeoutRef.current = setTimeout(() => {
      isLongPressRef.current = true;
      onLongPress();
    }, duration);
  }, [duration, onLongPress]);

  const handleMouseUp = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    if (!isLongPressRef.current) {
      onShortClick();
    }

    isLongPressRef.current = false;
  }, [onShortClick]);

  const handleMouseLeave = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
  }, []);

  return {
    onMouseDown: handleMouseDown,
    onMouseUp: handleMouseUp,
    onMouseLeave: handleMouseLeave,
    onTouchStart: handleMouseDown,
    onTouchEnd: handleMouseUp,
  };
}

export { useLongPress };
