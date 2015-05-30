package olduar

import "testing"

func BenchmarkInvetoryGet(b *testing.B) {
	//Create inventory
	inv := make(Inventory,0)
	
	//Add items
	for i := 0; i < 5000; i++ {
		inv.Add(&Item{
			Id: "Item",
		})
	}
	
	//Add new item to inventory
	shoe := ItemTemplateDirectory["useless_shoe"].GenerateItem()
	inv.Add(shoe)
		
	//Add more items
	for i := 0; i < 5000; i++ {
		inv.Add(&Item{
			Id: "Item",
		})
	}
	
	//Try to retrieve item
	inv.Get("useless_shoe")
}

func BenchmarkInvetoryRemove(b *testing.B) {
	//Create inventory & item
	inv := make(Inventory,0)
	
	//Add items
	for i := 0; i < 1000; i++ {
		inv.Add(&Item{
			Id: "Item",
		})
	}
	
	//Remove all items
	for len(inv) > 0 {
		item := inv[0]
		inv.Remove(item)
	}
}

func TestInventory(t *testing.T) {
	//Generate template
	uselessShoeTemplate := &ItemTemplate{
		Id: "useless_shoe",
		Name: "Useless shoe",
		Description: "Shoe full of holes that is completely useless",
	}
	ItemTemplateDirectory["useless_shoe"] = uselessShoeTemplate
	
	//Create inventory
	inv := make(Inventory,0)
	shoe := ItemTemplateDirectory["useless_shoe"].GenerateItem()
	
	//Add new item to inventory
	inv.Add(shoe)
	
	//Try to retrieve non-existing item
	if inv.Get("something_else") != nil {
		t.Error("Inventory returned something that is not there")
	}
	
	//Try to retrieve existing item
	if inv.Get("useless_shoe") != shoe  {
		t.Error("Inventory did not returned something that is there")
	}
	
	//Remove item from inventory
	inv.Remove(shoe)
	
	//Try to retrieve deleted item
	if inv.Get("useless_shoe") != nil {
		t.Error("Inventory did not returned something that is there")
	}
}