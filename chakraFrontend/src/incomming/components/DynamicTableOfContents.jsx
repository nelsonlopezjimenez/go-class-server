import { useEffect, useState } from 'react';

export default function DynamicTableOfContents({ sectionName, weekNumber }) {
  const [tocItems, setTocItems] = useState([]);
  const [activeId, setActiveId] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTOC();
  }, [sectionName, weekNumber]);

  useEffect(() => {
    // Wait a bit for content to load, then setup scroll tracking
    const timer = setTimeout(() => {
      setupScrollTracking();
    }, 1000);

    return () => clearTimeout(timer);
  }, [tocItems]);

  const fetchTOC = async () => {
    // Your existing fetch logic
    setLoading(true);
    try {
      const response = await fetch(
        `/api/sections/${sectionName}/week/${weekNumber}/toc`,
      );
      if (response.ok) {
        const data = await response.json();
        setTocItems(data.tocItems || []);
      } else {
        setTocItems(getDefaultTOCItems());
      }
    } catch (error) {
      setTocItems(getDefaultTOCItems());
    } finally {
      setLoading(false);
    }
  };

  const setupScrollTracking = () => {
    // First, ensure all headings have IDs
    const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
    headings.forEach((heading) => {
      if (!heading.id && heading.textContent) {
        heading.id = heading.textContent
          .toLowerCase()
          .replace(/[^\w\s-]/g, '')
          .replace(/\s+/g, '-')
          .trim();
      }
    });

    // Then setup scroll tracking
    const handleScroll = () => {
      const headingsWithIds = document.querySelectorAll(
        'h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]',
      );
      let currentActive = '';

      // Find the heading that's currently most visible
      headingsWithIds.forEach((heading) => {
        const rect = heading.getBoundingClientRect();
        if (rect.top <= 200 && rect.top >= -100) {
          currentActive = heading.id;
        }
      });

      if (currentActive !== activeId) {
        setActiveId(currentActive);
      }
    };

    // Throttle scroll events
    let ticking = false;
    const scrollListener = () => {
      if (!ticking) {
        requestAnimationFrame(() => {
          handleScroll();
          ticking = false;
        });
        ticking = true;
      }
    };

    window.addEventListener('scroll', scrollListener, { passive: true });
    handleScroll(); // Set initial state

    return () => {
      window.removeEventListener('scroll', scrollListener);
    };
  };

  const handleLinkClick = (e, itemId) => {
    e.preventDefault();
    const element = document.getElementById(itemId);
    if (element) {
      // Smooth scroll
      element.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
        inline: 'nearest',
      });

      // Update active state immediately
      setActiveId(itemId);

      // Update URL hash
      window.history.pushState(null, null, `#${itemId}`);
    }
  };

  const getDefaultTOCItems = () => [
    { id: 'learning-objectives', title: 'Learning Objectives', level: 2 },
    {
      id: 'course-structure-overview',
      title: 'Course Structure Overview',
      level: 2,
    },
    { id: 'daily-activities', title: 'Daily Activities', level: 2 },
  ];

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sticky top-4">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-200 rounded mb-4"></div>
          <div className="space-y-2">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-4 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <nav className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sticky top-4 max-h-96 space-y-0 overflow-y-auto">
      <h3 className="font-semibold mb-4 text-gray-900">
        Week {weekNumber} Contents
      </h3>

      <ul className="space-y-0.5">
        {tocItems.map((item, index) => (
          <li key={`${item.id}-${index}`} className="leading-tight">
            <a
              href={`#${item.id}`}
              onClick={(e) => handleLinkClick(e, item.id)}
              className={`block text-sm py-1 px-2 rounded-md transition-all duration-200 ${
                activeId === item.id
                  ? 'text-indigo-700 bg-indigo-50 border-l-4 border-indigo-500 font-medium'
                  : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-50 border-l-4 border-transparent'
              }`}
              style={{
                marginLeft: `${(item.level - 2) * 12}px`,
              }}
            >
              {item.title}
            </a>
          </li>
        ))}
      </ul>

      {/* Debug info - remove in production */}
      <div className="mt-4 pt-4 border-t border-gray-200 text-xs text-gray-500">
        Active: {activeId || 'none'}
      </div>
    </nav>
  );
}
