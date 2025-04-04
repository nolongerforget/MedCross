export function Card({ children, className }) {
    return (
      <div className={`bg-gray-800 p-6 rounded-xl shadow-md ${className}`}>
        {children}
      </div>
    );
  }
  
  export function CardContent({ children }) {
    return <div className="mt-4 text-white">{children}</div>;
  }
  