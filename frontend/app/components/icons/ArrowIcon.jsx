
  export default function ArrowIcon({
    className = '',
    direction = 'right',
  }) {
    return (
      <svg
        className={`${className} ${direction === 'left' ? 'flip' : ''} ${direction === 'down' ? 'rotate-90' : direction === 'up' ? 'rotate-[270deg]' : ''}`}
        width="24"
        height="24"
      >
        <path d="M7.293 4.707 14.586 12l-7.293 7.293 1.414 1.414L17.414 12 8.707 3.293 7.293 4.707z" />
      </svg>
    );
  }