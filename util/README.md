# Batch rename files

````
for myfile in ./*Slides*; 
do
    target=$(echo $myfile | sed -e 's#\(.*\) (Slides Attached)\(.*\)#\1\2#g')
    mv "$myfile" "$target"
done

for myfile in ./Lightning*; 
do
    target=$(echo $myfile | sed -e 's#\(.*\)Talk:\(.*\)#\1Talk\2#g')
    mv "$myfile" "$target"
done
````

Keynote CNCF Project Update - Liz Rice, Technology Evangelist, Aqua Security; Sugu Sougoumarane, CTO, PlanetScale Data; Colin Sullivan, Product Manager, Synadia Communications, Inc. & Andrew Jessup, Co-founder, Scytale Inc. - CNCF Project Update .pdf