import java.util.Hashtable;
import java.util.Map;

import java.io.File;

public class HashtableTaints
{
	public static void main(String args[])
	{
		if (args.length < 1) return;
		File f;

		Hashtable h1 = new Hashtable();
		h1.put("", args[0]);
		f = new File((String)h1.get("")); //if (Paths.get(h1).normalize().toString().startsWith("..

		Hashtable h2 = new Hashtable(h1);
		f = new File((String)h2.get("")); //CWEID 73
		f = new File((String)h2.elements().nextElement()); //CWEID 73
		f = new File(h2.toString()); //CWEID 73
		f = new File((String)h2.values().iterator().next()); //CWEID 73

		Map.Entry entry1 = (Map.Entry)h2.entrySet().iterator().next();
		f = new File((String)entry1.getValue()); //CWEID 73

		Hashtable h3 = new Hashtable();
		h3.put(args[0], "");
		f = new File((String)h3.keys().nextElement()); //CWEID 73
		f = new File((String)h3.keySet().iterator().next()); //CWEID 73
		f = new File(h3.toString()); //CWEID 73

		Map.Entry entry2 = (Map.Entry)h3.entrySet().iterator().next();
		f = new File((String)entry2.getKey()); //CWEID 73

		Hashtable h4 = new Hashtable();
		h4.putAll(h1);
		f = new File((String)h4.get("")); //CWEID 73
		f = new File((String)h4.elements().nextElement()); //CWEID 73

		Hashtable h5 = new Hashtable();
		h5.put("", "");
		Map.Entry entry3 = (Map.Entry)h5.entrySet().iterator().next();
		entry3.setValue(args[0]);
		f = new File((String)entry3.getValue()); //CWEID 73
		f = new File((String)h5.get("")); //CWEID 73
	}
}
