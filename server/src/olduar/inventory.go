package olduar

type Inventory []*Item

func (inv *Inventory) Remove(item *Item) {
	//TODO: Look for possible alternatives that could be faster
	newInventory := make(Inventory,len(*inv)-1)
	itemsAdded := 0
	for _, invItem := range *inv {
		if(invItem != item) {
			newInventory[itemsAdded] = invItem
			itemsAdded++
		}
	}
	*inv = newInventory
}

func (inv *Inventory) Add(item *Item) {
	*inv = append(*inv,item)
}

func (inv *Inventory) Get(entry string) *Item {
	for _, item := range *inv {
		if(item.Id == entry) {
			return item
		}
	}
	return nil
}
