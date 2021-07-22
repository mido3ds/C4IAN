import React, { useEffect, useState } from 'react';
import './Pagination.css';


function Pagination({ dataSize, paginate }) {
  const [pageNumber, setPageNumber] = useState(1);
  const [hasNext, setHasNext] = useState(pageNumber - Math.ceil(dataSize/6) < 0)

  useEffect(() => {
    let decButton = document.querySelector('#dec-btn');
    if (pageNumber > 1) {
      decButton.classList.replace('custom-btn-disabled', 'custom-btn');
    } else {
      decButton.classList.replace('custom-btn', 'custom-btn-disabled');
    }

    let incButton = document.querySelector('#inc-btn');
    if (hasNext) {
      incButton.classList.replace('custom-btn-disabled', 'custom-btn');
    } else {
      incButton.classList.replace('custom-btn', 'custom-btn-disabled');
    }
    
  }, [pageNumber, hasNext]);

  const btnClick = (inc) => {
    var paginationWrapper = document.querySelector('.pagination-wrapper');

    if (!inc && pageNumber > 1) {
      paginationWrapper.classList.add('transition-prev');
      setPageNumber(() => {
        paginate(pageNumber - 1);
        setHasNext((pageNumber - 1) - Math.ceil(dataSize/6) < 0)
        return pageNumber - 1
      });
    }
    
    else if (inc && hasNext) {
      paginationWrapper.classList.add('transition-next');
      setPageNumber(() => {
        paginate(pageNumber + 1);
        setHasNext((pageNumber + 1) - Math.ceil(dataSize/6) < 0)
        return pageNumber + 1
      });
    }

    setTimeout(cleanClasses, 500);
  };

  function cleanClasses() {
    var paginationWrapper = document.querySelector('.pagination-wrapper');

    if (paginationWrapper.classList.contains('transition-next')) {
      paginationWrapper.classList.remove('transition-next')
    } else if (paginationWrapper.classList.contains('transition-prev')) {
      paginationWrapper.classList.remove('transition-prev')
    }
  }

  return (
    <div className="pagination-wrapper">
      <svg id="dec-btn" className="custom-btn btn--prev" height="96" viewBox="0 0 24 24" width="96"
        onClick={btnClick.bind(this, false)} xmlns="http://www.w3.org/2000/svg">
        <path d="M15.41 16.09l-4.58-4.59 4.58-4.59L14 5.5l-6 6 6 6z" />
        <path d="M0-.5h24v24H0z" fill="none" />
      </svg>

      <div className="pagination-container">
        <div className="little-dot  little-dot--first"></div>
        <div className="little-dot">
          <div className="big-dot-container">
            <div className="big-dot"></div>
          </div>
        </div>
        <div className="little-dot  little-dot--last"></div>
      </div>

      <svg id="inc-btn" className="custom-btn btn--next" height="96" viewBox="0 0 24 24" width="96"
        onClick={btnClick.bind(this, true)} xmlns="http://www.w3.org/2000/svg">
        <path d="M8.59 16.34l4.58-4.59-4.58-4.59L10 5.75l6 6-6 6z" />
        <path d="M0-.25h24v24H0z" fill="none" />
      </svg>
    </div>
  );
};

export default Pagination;