import React, { useEffect, useState } from 'react';
import './Pagination.css';


function Pagination({ paginate, hasNext }) {
  const [pageNumber, setPageNumber] = useState(1);
  var paginationWrapper = document.querySelector('.pagination-wrapper');

  useEffect(() => {
    let incButton = document.querySelector('#inc-btn');
    if (hasNext) {
      incButton.classList.replace('custom-btn-disabled', 'custom-btn');
    } else {
      incButton.classList.replace('custom-btn', 'custom-btn-disabled');
    }
  }, [hasNext]);

  useEffect(() => {
    let decButton = document.querySelector('#dec-btn');
    if (pageNumber > 1) {
      decButton.classList.replace('custom-btn-disabled', 'custom-btn');
    } else {
      decButton.classList.replace('custom-btn', 'custom-btn-disabled');
    }
    
  }, [pageNumber]);

  const btnClick = (inc) => {
    if (!inc && pageNumber > 1) {
      paginationWrapper.classList.add('transition-prev');
      setPageNumber(pageNumber => {
        pageNumber = pageNumber - 1
        //paginate(pageNumber);
        return pageNumber;
      });
    }
    else if (inc && hasNext) {
      paginationWrapper.classList.add('transition-next');
      setPageNumber(pageNumber => {
        pageNumber = pageNumber + 1
        //paginate(pageNumber);
        return pageNumber;
      });
    }

    setTimeout(cleanClasses, 500);
  };

  function cleanClasses() {
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