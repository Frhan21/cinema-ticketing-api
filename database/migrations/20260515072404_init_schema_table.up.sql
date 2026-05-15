-- 1. Table user 
CREATE TABLE IF NOT EXISTS "users" {
    id uuid PRIMARY KEY DEFAULT gen_random_uuid()
    name VARCHAR(255) NOT NULL 
    email VARCHAR(255) NOT NULL UNIQUE 
    password VARCHAR(255) NOT NULL 
    role ENUM("admin", "user") DEFAULT "user"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
}

-- 2. Table studios
CREATE TABLE IF NOT EXISTS "studios" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    capacity INT NOT NULL, 
    facilites TEXT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
}

-- 3. Table Movies (Daftar film) 
CREATE TABLE IF NOT EXISTS "movies" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    genre VARCHAR(255) NOT NULL,
    duration INT, 
    description TEXT, 
    poster_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
}

-- 4. Table Schedules (Jadwal Film yang tayang)
CREATE TABLE IF NOT EXISTS "schedules" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL,
    studio_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL, 
    end_time TIMESTAMP NOT NULL, 
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- FOREIGN KEY constraints
    -- Movie key
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    -- Studio key
    FOREIGN KEY (studio_id) REFERENCES studios(id) ON DELETE CASCADE
}

-- 5. Table Seats (Kursi yang tersedia di studio)
CREATE TABLE IF NOT EXISTS "seats" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    studio_id UUID NOT NULL,
    seat_number VARCHAR(10) NOT NULL, 
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- FOREIGN KEY constraints
    FOREIGN KEY (studio_id) REFERENCES studios(id) ON DELETE CASCADE
}

-- 6. Table tickets (Tiket yang dibeli oleh user)
CREATE TABLE IF NOT EXISTS "tickets" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    schedule_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    status ENUM("pending", "paid", "cancelled") DEFAULT "pending",
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- FOREIGN KEY constraints
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,
    FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE CASCADE
}


-- 7. Table Transactions (Transaksi pembayaran tiket)
CREATE TABLE IF NOT EXISTS "transactions" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method ENUM(
        'cash',
        'qris',
        'transfer',
        'credit_card') NOT NULL,
    status ENUM("pending", "completed", "failed") DEFAULT "pending",
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- FOREIGN KEY constraints
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
}

-- 8. Table Transaction Items (Detail transaksi pembayaran tiket)
CREATE TABLE IF NOT EXISTS "transaction_items" {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    ticket_id UUID NOT NULL,

    -- FOREIGN KEY constraints
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
}