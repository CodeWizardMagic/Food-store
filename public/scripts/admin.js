document
  .getElementById("addProduct-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const productName = document.querySelector('input[name="productName"]').value;
    const productPrice = parseFloat(document.querySelector('input[name="productPrice"]').value);
    const productStock = parseInt(document.querySelector('input[name="productStock"]').value);

    try {
      const response = await fetch("/api/products", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          Name: productName,
          Price: productPrice,
          Stock: productStock,
        }),
      });

      if (response.ok) {
        document.getElementById("addProduct-form").reset();
        fetchProducts(1); 
      } else {
        const error = await response.text();
        alert(`Failed to add product: ${error}`);
      }
    } catch (error) {
      alert(`Failed to add product: ${error}`);
    }
  });

// Pagination state
let currentPage = 1;
const perPage = 5; 

async function fetchProducts(page = currentPage) {
  try {
    const url = new URL("/api/products", window.location.origin);
    url.searchParams.append("page", page);
    url.searchParams.append("per_page", perPage);

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch products: ${response.statusText}`);
    }
    const data = await response.json();

    const productsList = document.getElementById("products-list");
    productsList.innerHTML = "";

    data.products.forEach((product) => {
      const productItem = document.createElement("div");
      productItem.className = "main__products-item";
      productItem.innerHTML = `
        <div class="main__products-item-wrap">
          <div class="main__products-item-name">${product.Name}</div>
          <button class="main__products-item-edit" onclick="editProduct('${product.ID}', '${product.Name}', ${
            product.Price
          }, ${product.Stock})">✏️</button>
          <button class="main__products-item-delete" onclick="deleteProduct('${
            product.ID
          }')">❌</button>
        </div>
        <div class="main__products-item-price">${product.Price.toFixed(2)}₸</div>
        <div class="main__products-item-stock">Stock: ${product.Stock}</div>
      `;
      productsList.appendChild(productItem);
    });

    updatePagination(data.total, data.page, data.per_page);

  } catch (error) {
    console.error("Error fetching products:", error);
  }
}

function updatePagination(total, currentPage, perPage) {
  const totalPages = Math.ceil(total / perPage);
  const paginationDiv = document.getElementById("pagination");

  if (totalPages <= 1) {
    paginationDiv.style.display = "none";
    return;
  }

  paginationDiv.style.display = "flex";
  const prevButton = paginationDiv.querySelector(".main__pagination-button-prev");
  const nextButton = paginationDiv.querySelector(".main__pagination-button-next");
  const pageNumbersContainer = paginationDiv.querySelector(".main__pagination-numbers");

  pageNumbersContainer.innerHTML = "";

  prevButton.disabled = currentPage === 1;
  prevButton.onclick = () => {
    fetchProducts(currentPage - 1);
  };

  for (let i = 1; i <= totalPages; i++) {
    const pageBtn = document.createElement("button");
    pageBtn.textContent = i;
    pageBtn.className = "main__pagination-button main__pagination-button-page";
    pageBtn.disabled = i === currentPage;
    pageBtn.onclick = () => {
      fetchProducts(i);
    };
    pageNumbersContainer.appendChild(pageBtn);
  }

  nextButton.disabled = currentPage === totalPages;
  nextButton.onclick = () => {
    fetchProducts(currentPage + 1);
  };
}

async function deleteProduct(productId) {
  try {
    const response = await fetch(`/api/products/${productId}`, {
      method: "DELETE",
    });
    if (response.ok) {
      fetchProducts(currentPage); 
    } else {
      console.error("Error deleting product:", response.statusText);
    }
  } catch (error) {
    console.error("Error deleting product:", error);
  }
}

function editProduct(id, name, price, stock) {
  const productItem = document
    .querySelector(`button[onclick*="editProduct('${id}'"]`)
    .closest(".main__products-item");

  productItem.innerHTML = `
    <form class="main__products-item-update" onsubmit="submitUpdate(event, '${id}')">
      <input class="main__products-item-update-input" type="text" name="Name" value="${name}" required />
      <input class="main__products-item-update-input" type="number" name="Price" step="0.01" value="${price}" required />
      <input class="main__products-item-update-input" type="number" name="Stock" value="${stock}" required />
      <div class="main__products-item-update-buttons"> 
        <button class="main__products-item-update-buttons-save" type="submit">Save</button>
        <button class="main__products-item-update-buttons-cancel" type="button" onclick="fetchProducts(${currentPage})">Cancel</button>
      </div>
    </form>
  `;
}

async function submitUpdate(event, id) {
  event.preventDefault();

  const form = event.target;
  const updatedProduct = {
    Name: form.Name.value,
    Price: parseFloat(form.Price.value),
    Stock: parseInt(form.Stock.value),
  };

  try {
    const response = await fetch(`/api/products/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(updatedProduct),
    });

    if (response.ok) {
      fetchProducts(currentPage);
    } else {
      const error = await response.text();
      alert(`Failed to update product: ${error}`);
    }
  } catch (error) {
    alert(`Failed to update product: ${error}`);
  }
}

document
  .getElementById("filterProduct-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const minPrice = document.querySelector('input[name="minPrice"]').value;
    const maxPrice = document.querySelector('input[name="maxPrice"]').value;

    fetchProducts(1, { minPrice, maxPrice });
  });

async function fetchProducts(page = currentPage, filters = {}) {
  try {
    const url = new URL("/api/products", window.location.origin);
    url.searchParams.append("page", page);
    url.searchParams.append("per_page", perPage);

    if (filters.minPrice) url.searchParams.append("min_price", filters.minPrice);
    if (filters.maxPrice) url.searchParams.append("max_price", filters.maxPrice);

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch products: ${response.statusText}`);
    }
    const data = await response.json();

    const productsList = document.getElementById("products-list");
    productsList.innerHTML = "";

    data.products.forEach((product) => {
      const productItem = document.createElement("div");
      productItem.className = "main__products-item";
      productItem.innerHTML = `
        <div class="main__products-item-wrap">
          <div class="main__products-item-name">${product.Name}</div>
          <button class="main__products-item-edit" onclick="editProduct('${product.ID}', '${product.Name}', ${
            product.Price
          }, ${product.Stock})">✏️</button>
          <button class="main__products-item-delete" onclick="deleteProduct('${
            product.ID
          }')">❌</button>
        </div>
        <div class="main__products-item-price">${product.Price.toFixed(2)}₸</div>
        <div class="main__products-item-stock">Stock: ${product.Stock}</div>
      `;
      productsList.appendChild(productItem);
    });

    updatePagination(data.total, data.page, data.per_page);

  } catch (error) {
    console.error("Error fetching products:", error);
  }
}

window.onload = () => fetchProducts(1);