// package store

// import (
//     "errors"
//     "sync"
// )

// type StoreInfo struct {
//     StoreName string
//     AreaCode  string
// }

// type StoreManager struct {
//     stores map[string]StoreInfo
//     mu     sync.RWMutex
// }

// func NewStoreManager() *StoreManager {
//     return &StoreManager{
//         stores: map[string]StoreInfo{
//             "S00339218": {StoreName: "Store 1", AreaCode: "A1"},
//             "S01408764": {StoreName: "Store 2", AreaCode: "A2"},
//             // Add more stores as needed
//         },
//     }
// }

// func (sm *StoreManager) ValidateStore(storeID string) error {
//     sm.mu.RLock()
//     defer sm.mu.RUnlock()
    
//     if _, exists := sm.stores[storeID]; !exists {
//         return errors.New("store not found")
//     }
//     return nil
// }

package store

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "sync"
)

type StoreInfo struct {
    StoreName string
    AreaCode  string
}

type StoreManager struct {
    stores map[string]StoreInfo
    mu     sync.RWMutex
}

func NewStoreManager(csvPath string) (*StoreManager, error) {
    sm := &StoreManager{
        stores: make(map[string]StoreInfo),
    }

    // Open CSV file
    file, err := os.Open(csvPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create CSV reader
    reader := csv.NewReader(file)
    
    // Read header
    _, err = reader.Read()
    if err != nil {
        return nil, err
    }

    // Read and process each row
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        // CSV format: AreaCode, StoreName, StoreID
        sm.stores[record[2]] = StoreInfo{
            AreaCode:  record[0],
            StoreName: record[1],
        }
    }

    return sm, nil
}

func (sm *StoreManager) ValidateStore(storeID string) error {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    if _, exists := sm.stores[storeID]; !exists {
        return errors.New("store not found")
    }
    return nil
}

func (sm *StoreManager) GetStoreInfo(storeID string) (StoreInfo, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    if info, exists := sm.stores[storeID]; exists {
        return info, nil
    }
    return StoreInfo{}, errors.New("store not found")
}